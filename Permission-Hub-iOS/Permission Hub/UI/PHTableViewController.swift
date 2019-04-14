//
//  TransactionOverviewViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionVerificationStatus {
    case verified, unverified

    var color: UIColor {
        switch self {
        case .verified:
            return PHColors.topaz
        case .unverified:
            return PHColors.red
        }
    }
}

enum PHTableViewViewCellType {
    case notification(
        type: TransactionNotificationType,
        text: String)
    case warning(text: String)
    case description(image: UIImage?, date: Date, title: String, description: String)
    case plugin(image: UIImage?, text: String)
    case transactionItem(item: TransactionItem, isChecked: Bool)
    case selectionDisclosure(text: String)
    case selection(options: [String])
    case form(placeholder: String, text: String?, keyboardType: UIKeyboardType)
}

protocol PHTableViewControllerDelegate: class {
    func didSelect(item: PHTableViewViewCellType)
    func didDeselect(item: PHTableViewViewCellType)
}

struct ReloadedItems {

    let items: [PHTableViewViewCellType] = [
        .notification(
            type: .verification,
            text: "This company is verified"),
        .description(
            image: nil,
            date: Date(),
            title: "Personal details",
            description: "Please fill out your personal details."),
        .plugin(
            image: UIImage(named: "digid_button"),
            text: "Use the external DigiD plug-in to fill in your personal information (optional)."),
        .form(
            placeholder: "First name",
            text: "Gerard",
            keyboardType: .default),
        .form(
            placeholder: "Last name",
            text: "Huizinga",
            keyboardType: .default),
        .form(
            placeholder: "Date of birth",
            text: "04-11-1964",
            keyboardType: .numbersAndPunctuation),
        .form(
            placeholder: "Address",
            text: "Weesperplein 43-2, Amsterdam",
            keyboardType: .default),
        .form(
            placeholder: "Email",
            text: "gerard.huizinga@gmail.com",
            keyboardType: .emailAddress),
        .form(
            placeholder: "BSN number",
            text: "264036232",
            keyboardType: .numbersAndPunctuation)
    ]
}

class PHTableViewController: UIViewController {

    // MARK: - Private properties

    lazy var tableView: UITableView = {

        let tableView = UITableView()
        tableView.translatesAutoresizingMaskIntoConstraints = false

        tableView.rowHeight = UITableView.automaticDimension
        tableView.estimatedRowHeight = 66
        tableView.separatorStyle = .none

        tableView.dataSource = self
        tableView.delegate = self

        tableView.register(
            UITableViewCell.self,
            forCellReuseIdentifier: String(describing: UITableViewCell.self))

        tableView.register(
            TransactionOptionsTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionOptionsTableViewCell.self))

        tableView.register(
            TransactionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionTableViewCell.self))

        tableView.register(
            TransactionDescriptionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionDescriptionTableViewCell.self))

        tableView.register(
            TransactionPluginTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionPluginTableViewCell.self))

        tableView.register(
            TransactionNotificationTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionNotificationTableViewCell.self))

        tableView.register(
            FormTextInputCell.self,
            forCellReuseIdentifier: String(describing: FormTextInputCell.self))

        return tableView
    }()

    private var items: [PHTableViewViewCellType]

    private let activityIndicatorViewController: PHActivityIndicatorViewController = {

        let viewController = PHActivityIndicatorViewController()
        viewController.modalPresentationStyle = .overCurrentContext
        viewController.modalTransitionStyle = .crossDissolve

        return viewController
    }()

    // MARK: - Properties

    var shouldDisplayFooter: Bool {
        return false
    }

    weak var delegate: PHTableViewControllerDelegate?

    // MARK: - Initialization

    init(
        title: String? = nil,
        items: [PHTableViewViewCellType]) {

        self.items = items

        super.init(nibName: nil, bundle: nil)

        self.title = title
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = PHColors.lightGray

        // configure navigation bar
        configureNavigationBar()
        configureTableView()
        addHideKeyboardGestureRecognizer()
    }

    // MARK: - Configuration

    private func configureNavigationBar() {

        navigationItem.title = title
        navigationController?.navigationBar.tintColor = PHColors.greyishBrown
        navigationItem.backBarButtonItem = UIBarButtonItem(title: "", style: .plain, target: nil, action: nil)

        navigationController?.navigationBar.titleTextAttributes = [
            .font: PHFonts.wesBold(ofSize: 17),
            .foregroundColor: PHColors.greyishBrown
        ]
    }

    private func configureTableView() {

        view.addSubview(tableView)

        tableView.topAnchor.constraint(equalTo: view.topAnchor).isActive = true
        tableView.leftAnchor.constraint(equalTo: view.leftAnchor).isActive = true
        tableView.rightAnchor.constraint(equalTo: view.rightAnchor).isActive = true
        tableView.bottomAnchor.constraint(equalTo: view.bottomAnchor).isActive = true
    }

    private func addHideKeyboardGestureRecognizer() {

        let hideKeyboardTapGestureRecognizer = UITapGestureRecognizer(
            target: self,
            action: #selector(hideKeyboard))
        hideKeyboardTapGestureRecognizer.delegate = self
        view.addGestureRecognizer(hideKeyboardTapGestureRecognizer)
    }

    // MARK: - Selectors

    @objc private func hideKeyboard(_ sender: UITapGestureRecognizer) {
        view.endEditing(true)
    }
}

extension PHTableViewController: UITableViewDataSource {

    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return items.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {

        switch items[indexPath.row] {
        case .notification(let type, let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionNotificationTableViewCell.self),
                for: indexPath) as! TransactionNotificationTableViewCell

            cell.configure(withType: type, andText: text)

            return cell

        case .description(let image, let date, let title, let description):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionDescriptionTableViewCell.self),
                for: indexPath) as! TransactionDescriptionTableViewCell

            cell.configure(
                withImage: image,
                withDate: date,
                andTitle: title,
                andDescription: description)

            return cell

        case .transactionItem(let item, let isChecked):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: item.item,
                subtitle: (item.fields as NSArray).componentsJoined(by: ", "),
                shouldDisplayCheckmark: isChecked)
            cell.configure(withViewModel: viewModel)

            return cell

        case .warning(let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: text,
                subtitle: "")
            cell.configure(withViewModel: viewModel)

            return cell

        case .plugin(let image, let text):
            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionPluginTableViewCell.self),
                for: indexPath) as! TransactionPluginTableViewCell

            cell.configure(
                withImage: image,
                andText: text,
                callback: { [unowned self] in
                    self.present(self.activityIndicatorViewController, animated: false)
                    self.perform(
                        #selector(self.dismissActivityIndicator),
                        with: nil,
                        afterDelay: 2)
            })

            return cell

        case .selectionDisclosure(let text):
            let cell = UITableViewCell(style: .value1, reuseIdentifier: String(describing: UITableViewCell.self))
            cell.textLabel?.font = PHFonts.regular(ofSize: 14)
            cell.textLabel?.textColor = PHColors.greyishBrown
            cell.detailTextLabel?.text = "The Netherlands"
            cell.detailTextLabel?.font = PHFonts.regular(ofSize: 14)
            cell.detailTextLabel?.textColor = PHColors.greyishBrown.withAlphaComponent(0.5)
            cell.textLabel?.text = text
            cell.accessoryType = .disclosureIndicator

            return cell

        case .selection(let options):
            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionOptionsTableViewCell.self),
                for: indexPath) as! TransactionOptionsTableViewCell

            cell.configure(withOptions: options)

            return cell

        case .form(let placeholder, let text, let keyboardType):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: FormTextInputCell.self),
                for: indexPath) as! FormTextInputCell

            cell.configure(
                withPlaceholder: placeholder,
                andText: text,
                andKeyboardType: keyboardType) { text in
                    print(text)
            }

            return cell
        }
    }

    func tableView(_ tableView: UITableView, viewForFooterInSection section: Int) -> UIView? {

        let cell = UITableViewCell(style: .default, reuseIdentifier: nil)

        if shouldDisplayFooter {
            cell.textLabel?.font = PHFonts.regular(ofSize: 10)
            cell.textLabel?.textColor = PHColors.red
            cell.textLabel?.text = "* mandatory permissions"
        }

        return cell
    }

    // MARK: - Selectors

    @objc private func dismissActivityIndicator(_ sender: Any) {
        activityIndicatorViewController.dismiss(animated: true)
        
        items = ReloadedItems().items

        tableView.reloadData()
    }
}

extension PHTableViewController: UITableViewDelegate {

    func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
        delegate?.didSelect(item: items[indexPath.row])

        switch items[indexPath.row] {
        case .transactionItem:
            tableView.cellForRow(at: indexPath)?.isSelected = true

        default:
            break
        }
    }
}

extension PHTableViewController: UIGestureRecognizerDelegate {

    func gestureRecognizer(_ gestureRecognizer: UIGestureRecognizer, shouldReceive touch: UITouch) -> Bool {
        // add this to disable keyboard hiding tap gesture blocking tableview selection
        if touch.view!.isDescendant(of: tableView) {
            return false
        }

        return true
    }
}
